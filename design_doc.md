# Health Design Doc

## Things that need tracking
1. Calories
2. Weight
3. Fitness
   a. Strength
   b. Flexibility
   c. Endurances


## System Device

## Interface
cli based on bonzai

## Commands
- health add <type> <value> [<date>]
- health view

## Data Flow

```mermaid
graph TD
A[Manual Input]
B[Automatic Input]
C[File on Device]
D[Visualizations]

A --> C
B --> C
C --> D
```

## Input Types 
1. Manual Input 
   - User enters data directly into the app.
   - Examples: Weight, Calories consumed, Exercise details.
2. Automatic Input
   - Data is collected via connected devices or apps.
   - Examples: Fitness trackers, Smart scales.

## Architecture
each type can be defined as a module. 

Modules define the following 
- data input
- data storage (common location
- data visualization
- data instruction (next step)
