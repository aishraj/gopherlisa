## Things to do

- Rewrite the session storage manager
- Add a DB connection to mysql (needed for storing index) or to redis (depends)
- Fetch image metadata (as JSON)
- Resize and store tile images. (Each file would be identified by its hash ) and also the target image.
- Build the tile image index (sth like a lookup from color to hash)
  Table: (Hash could also be replaced by Instagram object id (subject to be verfied))
  |Color|Hash|filename|
-
