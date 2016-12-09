### Other antharun flags:


antharun --parameters --workflow

By default the antharun command uses a parameters file named parameters.json and a workflow file named workflow.json. 
If these files are named differently you’ll need to use the --parameters and/or --workflow flags to specify which files to use.

1.
To run the parameters found in this folder you'll need to run this:

antharun --parameters parameters.yml --workflow myamazingworkflow.json

_____________


antharun --inputPlateType

2. e.g. antharun --inputPlateType greiner384
This allows the type of input plate to be specified from the list of available Antha plate types. 
The available plates can be found by running the ```antharun lhplates``` command

 
_____________

antharun --inputPlates 

3. e.g. antharun --inputPlates inputplate.csv 
This allows user defined input plates to be defined. If this is not chosen antha will decide upon the layout.
More than one inputplate can be defined: this would be done like so:
antharun --inputPlates assemblyreagents.csv --inputPlates assemblyparts.csv

_____________

Config

4. An alternative to specifying plates as a flag is adding a Config section to the parameters file.
An input or output plate type can be specified by adding a config section to the parameters file as shown in configparameters.json

 "Config": {
        "InputPlateType": [
            "pcrplate_skirted_riser"
        ],
        "OutputPlateType": [
            "greiner384_riser"
        ]
    }
	
	
	
	
## Excercises

1. Check the available plates using ```antharun lhplates``` and change inputPlateType to one of the valid alternatives in the parameters file

## Next Steps
open [readme_LHComponents.md](readme_LHComponent.md) and continue