## antharun


Run this command from a folder containing your workflow.json file and parameters.yml file
as shown here. 
Running antharun without the ```--driver``` flag will use the manual driver.

________________

If a workflow or parameters set is changed you can rerun using antharun at any time. 

If you need to change the source code however, you'll need to recompile by running anthabuild.


## anthabuild

If you’ve added this alias this will build(recompile) all .an files in components into their corresponding .go files ready for execution. 
Whenever you change the source code of an antha element you must run anthabuild for the changes to take effect.

if you haven't set up the anthabuild alias you can do so by running this command:

#### Mac:

``` 
cat<<EOF>>$HOME/.bash_profile
alias anthabuild='make -C $HOME/go/src/github.com/antha-lang/elements'
EOF
```

#### Linux:

``` 
cat<<EOF>>$HOME/.bashrc
alias anthabuild='make -C $HOME/go/src/github.com/antha-lang/elements'
EOF
```

## Excercises

1. Modify the Sample.an file so that a additional Sample output is created called Sample2; modify the steps so that Sample2 is created in the same way as Sample, i.e. with the same Solution input and SolutionVolume.

You'll need to run anthabuild since the source code is being modified. If you get an error along the way, you'll need to resolve it before being able to run the modified element.

## Next Steps

Now go to [Folder B](../B_parallelruns/readme_drivers.md) to see how to run workflows with drivers.
