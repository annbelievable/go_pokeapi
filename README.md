BASIC DOCUMENTATION

Coded in golang version 1.13.8

How to use:
-- type in the pokemon name or id (the id uses number, eg: 1) and press enter
-- it will show basic information of the pokemon found with id, name, types, stats, encounter locations and methods
-- to stop the program, either close the terminal or use "control" + "c" on Windows or Linux, use "command" + "c" on Mac.




The API's website:
https://pokeapi.co/

The API's documentation:
https://pokeapi.co/docs/

To test the api, can use the website below:
https://reqbin.com/

All the APIs listed below are called using "GET" method:

--To get the id, name and type of the pokemon, use the API endpoint:
https://pokeapi.co/api/v2/pokemon/{id or name}/

The example is obtained from:
https://pokeapi.co/api/v2/pokemon/1/

--To get the location and encounter method of a pokemon, use the API endpoint:
https://pokeapi.co/api/v2/pokemon/{id or name}/encounters

The example is obtained from:
https://pokeapi.co/api/v2/pokemon/1/encounters

--To get the location, use the API endpoint:
https://pokeapi.co/api/v2/region/{id or name}/

The example is obtained from:
https://pokeapi.co/api/v2/region/kanto/

--To get the location, use the API endpoint:
https://pokeapi.co/api/v2/location/{id or name}/

The example is obtained from:
https://pokeapi.co/api/v2/location/67/

--To get the location area, use the API endpoint:
https://pokeapi.co/api/v2/location-area/{id or name}/

The example is obtained from:
https://pokeapi.co/api/v2/location-area/281

The structure of the region/location/area in the documentation.
Region is a larger place made up of many locations.
Location is a larger place made up of many location areas.
The structure:
Location region > Location > location area# go_pokeapi
