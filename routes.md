# Homework 4 Routes

| METHOD | PATH                         | FIELDS                           |
| :----: | ---------------------------- | -------------------------------- |
| GET    | /shard/my_id                 | NONE                             |
| GET    | /shard/all_ids               | NONE                             |
| GET    | /shard/members/{shard_id}    | NONE                             |
| GET    | /shard/count/{shard_id}      | NONE                             |
| GET    | /shard/changeShardNumber     | num={number}                     |
| GET    | /keyValue-store/{key}        | payload={payload}                |
| GET    | /keyValue-store/{key}        | val={value}<br>payload={payload} |
| DELETE | /keyValue-store/{key}        | payload={payload}                |
| GET    | /keyValue-store/search/{key} | payload={payload}                |
| GET    | /view                        | NONE                             |
| PUT    | /view                        | ip_port={NewIPPort}              |
| DELETE | /view                        | ip_port={RemovedIPPort}          |


The endpoints from the previous assignments are back as well. Notice that there has been a slight change to GET  /keyValue-store/&lt;key> and GET /keyValue-store/search/&lt;key>, in that you should also return the shard id of the node who is holding that particular key.

Additional logic will also be needed in the adding and removing of nodes to the view to keep the number of nodes in each shard about the same and to maintain fault tolerance.

For example, imagine you had 5 nodes split into 2 shards like so: [A,B,C],[D,E]. Now imagine I removed E, we would be left with: [A,B,C],[D]. In this case we would want to move one of the other nodes into the second shard, like so: [A,B],[D,C]

Alternatively, if we had 6 nodes split into 3 shards: [A,B],[C,D],[E,F] deleting one node would leave us with: [A,B],[C,D],[E]. In this case the only thing we can do is to remove a shard, like so: [A,B],[C,D,E]

 

# Route Details

## /shard/my_id
Should return the container's shard id
```json
{
  "id": "<container’s ShardID>"
}
```
```js
// Status: 200
// Method: GET
// Fields: None
```
## /shard/all_ids
Should return a list of all shard ids in the system as a string of comma separated values
```json
{
  "result": "Success",
  "shard_id": "0,1,2"
}
```
```js
// Status: 200
// Method: GET
// Fields: None
```
## /shard/members/{shard_id}
Should return a list of all members in the shard with id &lt;shard_id>. Each member should be represented as an ip-port address. (Again, the same one you pass into VIEW)
```json
{
  "result": "Success",
  "members": "176.32.164.2:8080,176.32.164.3:8080”“176.32.164.2:8080,176.32.164.3:8080"
}
```
```js
// Status: 200
// Method: GET
// Fields: None
```
If the &lt;shard_id> is invalid, please return:
```json
{
"result": "Error",
"msg": "No shard with id <shard_id>"
}
```
```js
// Status: 404
```
## /shard/count/{shard_id}
Should return the number of key-value pairs that shard is responsible for as an integer
```json
{
    "result": "Success",
    "Count": "<numberOfKeys>",
}
```
```js
// Status: 200
// Method: GET
// Fields: None
```
If the &lt;shard_id> is invalid, please return:
```json
{
    "result": "Error",
    "msg": "No shard with id <shard_id>",
}
```
```js
// Status: 404
```
## /shard/changeShardNumber
Should initiate a change in the replica groups such that the key-values are redivided across &lt;number> groups and returns a list of all shard ids, as in GET /shard/all_ids

```json
{
    "result": "Success",
    "shard_ids": "0,1,2",
 }

```js
// Status: 200
```
## /keyValue-store/{key}
```js
// Method: GET
// Status: 200
```
## /keyValue-store/{key}
```js
// Method: GET
// Status: 200
```
## /keyValue-store/{key}
```js
// Method: DELETE
// Status: 200
```
## /keyValue-store/search/{key}
```js
// Method: GET
// Status: 200
```
## /view
```js
// Method: GET
// Status: 200
```
## /view
```js
// Method: PUT
// Status: 200
```
## /view
```js
// Method: DELETE
// Status: 200
```