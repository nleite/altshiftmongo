Using MongoDB Arrays with GOLANG
================================

This post was inspired by this other blog post from `OpsDash`_ where the author
explains, in a very practical way, how to make use of `SQL99 arrays type`_ in
PostgreSQL to support faster queries.

The situation
-------------

Let's see how that would be looking like using MongoDB instead.

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

Further Reading
---------------


Other Topics
------------

.. _`OpsDash`: https://www.opsdash.com/blog/postgres-arrays-golang.html?h=1
.. _`SQL99 arrays type`: https://www.iso.org/standard/26197.html
