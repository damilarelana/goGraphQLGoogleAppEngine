### goGraphQL Server

A golang implementation of GraphQL server, via the Google App Engine server-less platform

***

* Golang
* GraphQL
* Google App Engine


***

The GraphQL server is hosted on the Google App Engine at `https://graphqlserver-259904.appspot.com/graphql`. To test the server, the following mutation and query examples would be used.

#### *Mutations*
To create a user `Banner`, run `https://graphqlserver-259904.appspot.com/graphql?query=mutation{createUser(name:"Banner"){id}}` as a `GET` request in [Postman](https://www.getpostman.com/downloads/).

To create users (`John`, `Mark`, `Bob`), run `https://graphqlserver-259904.appspot.com/graphql?query=mutation{john:createUser(name:"John"){id},bob:createUser(name:"Bob"){id},mark:createUser(name:"Mark"){id}}` as a `GET` request in [Postman](https://www.getpostman.com/downloads/).

To create posts, run `https://graphqlserver-259904.appspot.com/graphql?query=mutation{a:createPost(userID:"5768037999312896",content:"Hi!"){id,content},b:createPost(userID:"5768037999312896",content:"lol"){id,content},c:createPost(userID:"5768037999312896",content:"GraphQL is pretty cool!"){id,content}}` as a `GET` request in [Postman](https://www.getpostman.com/downloads/).


#### Queries

To query posts, run `https://graphqlserver-259904.appspot.com/graphql?query={posts{totalCount,nodes{id,content,createdAt}}}` as a `GET` request in [Postman](https://www.getpostman.com/downloads/).

To query users, run `https://graphqlserver-259904.appspot.com/graphql?query={user(id:"5646874153320448"){name,posts{totalCount,nodes{content}}}}` as a `GET` request in [Postman](https://www.getpostman.com/downloads/).