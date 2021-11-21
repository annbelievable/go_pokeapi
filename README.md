# Golang POKEAPI
This is a simple terminal application that is used to search for pokemon data.
This project is coded in golang version 1.13.8 and using the API, PokeAPI.
For more detailed information about the API used in this project, you can visit [PokeAPI](https://pokeapi.co/).
The API's documentation page can be found here [PokeAPI Documentation](https://pokeapi.co/docs/).
This project uses the library Cobra to build the terminal command part, the website 
and repoitory for the library can be found at [Cobra website](https://cobra.dev/) and [Cobra Github](github.com/spf13/cobra) respectively.

## How to use
- Make sure your environment is installed with golang with version 1.13.8 or above.
- Git clone the repo to a folder.
- In that cloned folder, run the command `go install pokeapi` to make sure all the required dependencies are installed.
- Also make sure your GOPATH is being setup correctly.
- You can now use the command `pokeapi pokemon -q {query}` to search for pokemon data.
- The `{query}` can be either the pokemon name or pokemon id number.
- If youre searching by name please do not include any space, any string after the first space will be ignored. 
- If there is result, it will show basic information of the pokemon found with id, name, types, stats, encounter locations and methods.
- If there is no result, it will tell you that there is an error and no result is found.

## License
There is no license for this project at the moment.