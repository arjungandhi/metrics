package strava

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"

	"github.com/arjungandhi/metrics/pkg/metric"
	"github.com/arjungandhi/metrics/pkg/store"
	"golang.org/x/oauth2"
)

const configKey = "provider.strava"

var oauthConfig = &oauth2.Config{
	Endpoint: oauth2.Endpoint{
		AuthURL:  "https://www.strava.com/oauth/authorize",
		TokenURL: "https://www.strava.com/oauth/token",
	},
	RedirectURL: "http://localhost:8089/callback",
	Scopes:      []string{"activity:read_all"},
}

// config holds the persisted settings for the Strava provider.
type config struct {
	ClientID     string       `json:"client_id"`
	ClientSecret string       `json:"client_secret"`
	Token        *oauth2.Token `json:"token"`
	LastSync     *time.Time   `json:"last_sync,omitempty"`
}

// Provider implements provider.Provider for Strava.
type Provider struct{}

func New() *Provider {
	return &Provider{}
}

func (p *Provider) Name() string { return "strava" }

// Setup walks the user through OAuth2 setup and stores credentials.
func (p *Provider) Setup(s store.Store) error {
	var clientID, clientSecret string
	fmt.Println("Strava API credentials (from https://www.strava.com/settings/api)")
	fmt.Print("Client ID: ")
	if _, err := fmt.Scanln(&clientID); err != nil {
		return err
	}
	fmt.Print("Client Secret: ")
	if _, err := fmt.Scanln(&clientSecret); err != nil {
		return err
	}

	oauthConfig.ClientID = clientID
	oauthConfig.ClientSecret = clientSecret

	// Start local server to catch the OAuth callback.
	ln, err := net.Listen("tcp", "localhost:8089")
	if err != nil {
		return fmt.Errorf("starting callback server: %w", err)
	}
	defer ln.Close()

	tokenCh := make(chan *oauth2.Token, 1)
	errCh := make(chan error, 1)

	mux := http.NewServeMux()
	mux.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		if code == "" {
			errCh <- fmt.Errorf("no code in callback")
			fmt.Fprintln(w, "Error: no code received.")
			return
		}
		tok, err := oauthConfig.Exchange(context.Background(), code)
		if err != nil {
			errCh <- err
			fmt.Fprintf(w, "Error: %v", err)
			return
		}
		tokenCh <- tok
		fmt.Fprintln(w, "Success! You can close this window.")
	})

	srv := &http.Server{Handler: mux}
	go srv.Serve(ln)

	authURL := oauthConfig.AuthCodeURL("state", oauth2.AccessTypeOffline)
	fmt.Printf("\nOpen this URL in your browser:\n  %s\n\nWaiting for authorization...\n", authURL)

	var tok *oauth2.Token
	select {
	case tok = <-tokenCh:
	case err := <-errCh:
		return err
	}

	srv.Shutdown(context.Background())

	cfg := config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Token:        tok,
	}
	return saveConfig(s, &cfg)
}

// Sync fetches activities from Strava and writes metrics to the store.
func (p *Provider) Sync(s store.Store) error {
	cfg, err := loadConfig(s)
	if err != nil {
		return fmt.Errorf("load strava config (run setup first): %w", err)
	}

	oauthConfig.ClientID = cfg.ClientID
	oauthConfig.ClientSecret = cfg.ClientSecret

	// Use oauth2 token source for automatic refresh.
	ts := oauthConfig.TokenSource(context.Background(), cfg.Token)
	httpClient := oauth2.NewClient(context.Background(), ts)

	// Fetch activities since last sync, or all activities on first run.
	var after *time.Time
	if cfg.LastSync != nil {
		after = cfg.LastSync
	}

	activities, err := fetchActivities(httpClient, after)
	if err != nil {
		return err
	}

	for _, a := range activities {
		labels := map[string]string{
			"type": a.Type,
			"name": a.Name,
		}
		ts := a.StartDate

		if err := s.AddDataPoint("strava.distance", metric.DataPoint{
			Time: ts, Value: a.Distance, Labels: labels,
		}); err != nil {
			return err
		}
		if err := s.AddDataPoint("strava.elevation", metric.DataPoint{
			Time: ts, Value: a.TotalElevationGain, Labels: labels,
		}); err != nil {
			return err
		}
		if err := s.AddDataPoint("strava.moving_time", metric.DataPoint{
			Time: ts, Value: float64(a.MovingTime), Labels: labels,
		}); err != nil {
			return err
		}
		if err := s.AddDataPoint("strava.average_speed", metric.DataPoint{
			Time: ts, Value: a.AverageSpeed, Labels: labels,
		}); err != nil {
			return err
		}
		if err := s.AddDataPoint("strava.max_speed", metric.DataPoint{
			Time: ts, Value: a.MaxSpeed, Labels: labels,
		}); err != nil {
			return err
		}
		if err := s.AddDataPoint("strava.activities", metric.DataPoint{
			Time: ts, Value: 1, Labels: labels,
		}); err != nil {
			return err
		}
	}

	// Persist refreshed token and last sync time.
	newTok, err := ts.Token()
	if err == nil {
		cfg.Token = newTok
	}
	now := time.Now()
	cfg.LastSync = &now
	if err := saveConfig(s, cfg); err != nil {
		return fmt.Errorf("saving config after sync: %w", err)
	}

	fmt.Printf("Synced %d activities from Strava\n", len(activities))
	return nil
}

// activity is a subset of the Strava activity response.
type activity struct {
	Name               string    `json:"name"`
	Type               string    `json:"type"`
	StartDate          time.Time `json:"start_date"`
	Distance           float64   `json:"distance"`
	MovingTime         int       `json:"moving_time"`
	TotalElevationGain float64   `json:"total_elevation_gain"`
	AverageSpeed       float64   `json:"average_speed"`
	MaxSpeed           float64   `json:"max_speed"`
}

func fetchActivities(client *http.Client, after *time.Time) ([]activity, error) {
	var all []activity
	page := 1

	for {
		url := fmt.Sprintf(
			"https://www.strava.com/api/v3/athlete/activities?page=%d&per_page=100",
			page,
		)
		if after != nil {
			url += fmt.Sprintf("&after=%d", after.Unix())
		}

		resp, err := client.Get(url)
		if err != nil {
			return nil, fmt.Errorf("fetching activities: %w", err)
		}

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			resp.Body.Close()
			return nil, fmt.Errorf("strava API error %d: %s", resp.StatusCode, body)
		}

		var batch []activity
		err = json.NewDecoder(resp.Body).Decode(&batch)
		resp.Body.Close()
		if err != nil {
			return nil, fmt.Errorf("decoding activities: %w", err)
		}

		if len(batch) == 0 {
			break
		}

		all = append(all, batch...)
		page++
	}

	return all, nil
}

func loadConfig(s store.Store) (*config, error) {
	raw, err := s.GetConfig(configKey)
	if err != nil {
		return nil, err
	}
	var cfg config
	if err := json.Unmarshal(raw, &cfg); err != nil {
		return nil, fmt.Errorf("parsing strava config: %w", err)
	}
	return &cfg, nil
}

func saveConfig(s store.Store, cfg *config) error {
	data, err := json.Marshal(cfg)
	if err != nil {
		return err
	}
	return s.SetConfig(configKey, data)
}
