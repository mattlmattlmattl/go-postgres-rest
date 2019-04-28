# go-postgres-rest
## Very simple REST api implemented with Go and PostgresQL.

Uses gorilla/mux and pq.  Inputs and emits JSON (other than a lone id).

The database is a single table, "item", with three fields, "id" (auto generated), "name" and "notes".

To build the server: 
```
# go build
```
(should create go-postgres-rest executable)

To run the server: 
```
# ./go-postgres-rest
```

To view a list of all items, go to 
http://localhost:8000/items

To see single item with id=4: 
http://localhost:8000/item/4

To add an item
```
# curl --data "{\"name\":\"Tester\",\"notes\":\"testeroni\"}" http://localhost:8000/items
```

To update item with id=9
```
# curl -X PUT -d "{\"name\":\"FatCat\"}" -d "{\"notes\":\"blah\"}"  http://localhost:8000/item/9
```

To delete item with id=7
```
# curl -X DELETE http://localhost:8000/item/7
```
