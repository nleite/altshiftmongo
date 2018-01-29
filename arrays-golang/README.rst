Using MongoDB Arrays with GOLANG
================================

This post was inspired by this other blog post from `OpsDash`_ where the author
explains, in a very practical way, how to make use of `SQL99 arrays type`_ in
PostgreSQL to support faster queries.

The situation
-------------

Let's see how that would be looking like using MongoDB.

**Array Types in BSON**

BSON, which is the underlying MongoDB data format, supports array types::

  Array - The document for an array is a normal BSON document with integer
  values for the keys, starting with 0 and continuing sequentially. For example,
  the array ['red', 'blue'] would be encoded as the document
  {'0': 'red', '1': 'blue'}. The keys must be in ascending numerical order.

Which basically translates into the following::

  BSON Arrays are pretty darn flexible.

No specific array type needs to be defined for that field, and the array
elements can be of any valid BSON type.

We can define arrays of integer elements:

.. code-block:: js

  db.arrays.insert( { int_array : [1,2,3]})
  db.arrays.find( { int_array : { "$type" : "array" } } )
  {
    "_id" : ObjectId("5a36a520f16412b7a5095099"),
    "some_array" : [ 1, 2, 3 ]
  }

Or string elements:

.. code-block:: js

  db.arrays.find( {string_array: { '$type': 'array' }} )
  {
    "_id" : ObjectId("5a36a60cf16412b7a509509a"),
    "string_array" : [ "hello", "world" ]
  }


Or even arrays with mixed types:

.. code-block:: js

  db.arrays.insert( { any_array: [ "hello", 1, 2, "world" ] } )
  db.arrays.find( {any_array: { '$type': 'array' }} )
  {
    "_id" : ObjectId("5a36a6a4f16412b7a509509b"),
    "any_array" : [ "hello", 1, 2, "world" ]
  }

Given this flexibility, of allowing multiple elements of different types in an
array, there are some nice operators that allows us to perform a few interesting
queries.
If we want to ask the database for:
*all documents where the second element of the array is an number*, we can:

.. code-block:: js

  db.arrays.find( { 'any_array.1': {'$type': 'number'} })
  {
    "_id" : ObjectId("5a36a6a4f16412b7a509509b"),
    "any_array" : [ "hello", 1, 2, "world" ]
  }


In MongoDB we don't need to explicitly define a data type for any given key in a
document. However, we can define `document validation`_ rules, to ensure that
any given array, or field by that matter, is of a given set of types.

But I digress, what is realy handy about this is the fact that we can move data
arround and settle on a particular data type, or just use a few, within the
array elements once we are certain about what those elements types should be.

Operations with arrays
----------------------

When using arrays, you might want to use them in full throotle. What does that mean?
Well, that means that you might want to:

* Match documents on single elements of an array by their value
* Match documents on element array field position value
* Update individual elements of an array when matching a update selector
* Make use of array operators (min, max, slice, push, pull ... )
* Update all elements of a given array that match a filter

And you can absolutely do that with MongoDB.

**Match documents on single elements of an array by their value**

As simple as expressing a query in MongoDB:

.. code-block:: js

  db.arrays.find({ string_array: "hello"})
  {
    "_id" : ObjectId("5a36a60cf16412b7a509509a"),
    "string_array" : [ "hello", "world" ]
  }


**Match documents on element array index position value**

In this case we will make use of the `dot notation`_ and define the array index:

.. code-block:: js

  db.arrays.find( { "int_array.1": { "$gte" : 2 } } )
  {
    "_id" : ObjectId("5a37168049046afc0b63c7c2"),
    "int_array" : [ 1, 2, 3 ]
  }


**Match documents on subdocument field array elements**

Things get a lot more interesting when dealing with subdocuments as array
elements.

.. code-block:: js

  db.arrays.insert({
    "_id": ObjectId("5a37168049046afc0b63c7c3"),
    "complex_array": [
      {
        "name": "Bernie",
        "grade": 10,
        "city": "New York"
      },
      {
        "name": "Ernie",
        "grade": 12,
        "city": "New York"
      },
      {
        "name": "Dottie",
        "grade": 10,
        "city": "Porto"
      }
    ]
  })

Given this array with subdocuments, we can match the document based on any field
of the inner array subdocument fields:

.. code-block:: js

  db.arrays.find({ "comples_array.name": "Bernie"  })

If our query is looking for the composition of more than one field in a
subdocument, we will have to use ``$elemMatch`` operator:

.. code-block:: js

  db.arrays.find({
    "comples_array": {
      "$elemMatch": {
        "name": "Dottie",
        "city": "Porto"  }
      }
    })


Using GO
--------

There are a few opensource community supported `MongoDB GO`_ libraries out
there.

The most popular GO MongoDB library (at the time of this writting) is
`mgo`_.
Given the popular vote, we will be using **mgo** in this set of examples.

The first thing we need to do is establish a database connection/session

.. code-block:: go

  package main

  import (
  	"log"

  	"gopkg.in/mgo.v2"
  )

  func main() {
  	session, err := mgo.Dial("localhost:27017")
  	if err != nil {
  		log.Fatal("Make sure your ``mongod`` is up and running on port **27017**")
  		panic(err)
  	}
  	defer session.Close()

  	// ... now we are ready to start using our db
  }

Once we have session, we can then initialize a **DB** and **Collection**
instances.

.. code-block:: go

  db := session.DB("altshiftmongo")
  collection := db.C("arrays")

We will use collections to store documents.

.. code-block:: go

  type SomeArray struct {
    SomeArray []int `json:"some_array" bson:"some_array"`
  }

Let's start by using a simple ``SomeArray`` struct to operate data in go.

.. code-block:: go

  func main(){
    // ...
    doc1 := SomeArray{[]int{1, 2, 3, 4}}
    insertErr := collection.Insert(&doc1)
    if insertErr != nil {
      log.Fatal(err)
    }

  }

To store documents we can simply pass pointer to the structure we want store,
to the ``Insert`` function of our ``Collection`` object.

`mgo`_ will be marshalling our ``SomeArray`` type into a ``BSON`` object and
sent it to MongoDB. All communication between the client application and
MongoDB are done by exchanging `wire protocol`_ messages, which are
themselves BSON based messages.

.. code-block:: sh

  mongo altshiftmongo --eval 'db.arrays.find().pretty()' --quiet
  {
    "_id" : ObjectId("5a6c918ab0e288905514ada8"),
      "some_array" : [
        1,
        2,
        3,
        4
      ]
  }

If go look into the database using the ``mongo`` shell client, we can see our
newly inserted document.

Collection with different documents
------------------------------------

As I mentioned in the first section of this post, we can infact have different
*shapes* of documents throughout our collections.

.. code-block:: go

  type StringArray struct {
    StringArray []string `json:"string_array" bson:"string_array"`
  }

In this example, apart from the ``SomeArray`` type, there's another type called
``StringArray`` that, as name indicates, a list of string type values.

.. code-block:: go

  doc2 := StringArray{[]string{"bernie", "ernie", "dottie"}}
  insertErr := collection.Insert(&doc2)
  if insertErr != nil {
    log.Fatal(err)
  }

Looking back into the database we can find the two different documents in the
same collection.

.. code-block:: sh

  mongo altshiftmongo --eval 'db.arrays.find().pretty()' --quiet
  {
  "_id" : ObjectId("5a6c9551b0e288905514adcc"),
  "some_array" : [
    1,
    2,
    3,
    4
  ]
  }
  {
  "_id" : ObjectId("5a6c9551b0e288905514adcf"),
  "string_array" : [
    "bernie",
    "ernie",
    "dottie"
  ]
  }


Find Documents
--------------

Obviously, the driver needs to allow a full set of CRUD operations, which include
the capability of retrieving back to the client application, documents based on some query / filtering criteria.

In traditional SQL we would be using a ``SELECT`` statement, in MongoDB we can
simply express query based on a struct expressing the matching criteria.




Further Reading
---------------

MongoDB offers a wide variety of `array query`_ and `array update`_ operators.
That allows us to, from the database prespective, have a quite extensive
manuverability on how to deal with array fields.


Other Topics
------------

.. _`mgo`: https://labix.org/mgo
.. _`OpsDash`: https://www.opsdash.com/blog/postgres-arrays-golang.html?h=1
.. _`SQL99 arrays type`: https://www.iso.org/standard/26197.html
.. _`array query`: https://docs.mongodb.com/manual/reference/operator/query-array/
.. _`array update`: https://docs.mongodb.com/manual/reference/operator/update-array/
.. _`MongoDB GO libraries`: https://docs.mongodb.com/ecosystem/drivers/community-supported-drivers/
