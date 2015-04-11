## Specification

The API aims to be simple to use and understand. A full range of methods are supported to get, create and update avatar images.

### GET

`/:hash[/:backup][/:size]`

Fetch an avatar, optionally at a specific size.

#### Parameters

* `hash`: SHA1 of a unique identifier (i.e. email, id, etc)
* `backup`: SHA1 to use if the given `:hash` does not exist
* `size`: one of original, large, medium, or small:
    * original: 1024x
    * large: 512x
    * medium: 256x
    * small: 128x

When the provided `:hash` does not exist, and a `:backup` is provided, the backup is treated as the requested hash.

When the provided `:size` is an integer, the closest **larger** size will be used. For example, if the size was `300`, the size will be converted to `large`. If no `:size` is provided, it defaults to `medium`.

#### Response Headers

* `Location: (Image File URL)`

#### Response Status

* `302`: redirect to image file

_The result of this call will **never** return a 404! If the requested size does not exist, return the best available size instead._

### HEAD

`/:hash`

Check if an avatar exists.

Response status:

* `404`: not found
* `204`: success

### POST

`/:hash`

Image is processed into a square and sizes are immediately created and stored on S3.

#### Parameters

* `avatar`: image file upload in the post body
* `token`: a [JWT](http://jwt.io/) containing: exp, hash

#### Request Headers

* `Content-Type: multipart/form-data`

#### Response Status

* `201`: success

#### Example Response

```json
{
  "data": {
    "type": "jpg",
    "best": "medium",
    "sizes": {
      "small": {
        "href": "\/4e1243bd22c66e76c2ba9eddc1f91394e57f9f83\/small",
        "path": "4\/e1\/4e1243bd22c66e76c2ba9eddc1f91394e57f9f83.small.jpg",
        "url": "\/\/s3.amazonaws.com\/s3-bucket.example.com\/4\/e1\/4e1243bd22c66e76c2ba9eddc1f91394e57f9f83.small.jpg"
      },
      "medium": {
        "href": "\/4e1243bd22c66e76c2ba9eddc1f91394e57f9f83\/medium",
        "path": "4\/e1\/4e1243bd22c66e76c2ba9eddc1f91394e57f9f83.medium.jpg",
        "url": "\/\/s3.amazonaws.com\/s3-bucket.example.com\/4\/e1\/4e1243bd22c66e76c2ba9eddc1f91394e57f9f83.medium.jpg"
      }
    }
  },
  "error": false
}
```


### DELETE

`/:hash`

Delete the given avatar and purge all existing sizes.

#### Parameters

* `token`: a [JWT](http://jwt.io/) containing: exp, hash

#### Response Status

* `404`: not found
* `504`: failed to delete some images
* `204`: success
