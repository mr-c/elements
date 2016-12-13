### Other antharun flags:
 
A full list of optional flags which can be used with antharun is available by running ```antharun --help```

### antharun --parameters --workflow

By default the antharun command uses a parameters file named parameters.json and a workflow file named workflow.json. 
If these files are named differently you’ll need to use the --parameters and/or --workflow flags to specify which files to use.


To run the parameters found in this folder you'll need to run this:

```bash
antharun --parameters parameters.yml --workflow myamazingworkflow.json
```

### antharun --bundle 

By default the antharun command uses a parameters file named parameters.json and a workflow file named workflow.json. 
If these files are named differently you’ll need to use the --parameters and/or --workflow flags to specify which files to use.


To run the combined parameters and workflow bundle found in this folder you'll need to run this:


```bash
antharun --bundle bundle.json 
```


### antharun --inputPlateType

e.g. 
```bash
antharun --inputPlateType greiner384
```

This allows the type of input plate to be specified from the list of available Antha plate types. 
The available plates can be found by running the ```antharun list plates``` command

 

### antharun --inputPlates 

e.g. 
```bash
antharun --inputPlates inputplate.csv 
```

This allows user defined input plates to be defined. If this is not chosen antha will decide upon the layout.
More than one inputplate can be defined: this would be done like so:

```bash
antharun --inputPlates assemblyreagents.csv --inputPlates assemblyparts.csv
```


### Config

4. An alternative to specifying plates as a flag is adding a Config section to the parameters file.
A series of desired input or output plates (in order of preference) can be specified by adding a config section to the parameters file as shown in configparameters.json

 "Config": {
        "InputPlateType": [
            "pcrplate_skirted_riser"
        ],
        "OutputPlateType": [
            "greiner384_riser"
        ]
    }
	
There are many other preferences which can be specified in the config, such as tip position preferences to whether you want Antha to compensate for evaporation. 	
	
## Excercises

1. Check the available plates using ```antharun list plates``` and change inputPlateType to one of the valid alternatives in the parameters file

## Next Steps
open [readme_LHComponents.md](readme_LHComponent.md) and continue
