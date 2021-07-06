## `CS5429 Distributed Computing: Lab Practical 01`

`We can access client via api endpoints and client will execute RPC in server and return the outputs in response`

`server will run in port 9000 in localhost and client will run in port 7000 in localhost`

##### `Receive a list of all available vegetables and display.`
`client will expose http://localhost:7000/vegetables/get/all endpoint to excute that function`

##### `Get the price per kg of a given vegetable and display`
`client will expose http://localhost:7000/vegetables/get/pricePerKg?Name=Carrot endpoint to excute that function`

##### `Get the available amount of kg of a given vegetable and display`
`client will expose http://localhost:7000/vegetables/get/availableAmount?Name=Carrot endpoint to excute that function`

##### `Send a new vegetable name to the server to be added to the server file`
`client will expose http://localhost:7000/vegetables/add?Name=Cabbage&PricePerKg=430&AvailableAmountOfKg=20 endpoint to excute that function`

##### `Send new price or available amount for a given vegetable to be updated in the server file.`
`client will expose http://localhost:7000/vegetables/update?Name=Carrot&PricePerKg=430&AvailableAmountOfKg=20 endpoint to excute that function`