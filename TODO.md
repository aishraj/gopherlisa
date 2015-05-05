## Things to do

- Build JSON types
- Fetch image metadata (as JSON), that too by traversing through the pagination.
- Resize and store tile images. (Each file would be identified by its hash ) and also the target image.
- Add a DB connection to mysql (needed for storing index) or to redis (depends)
- Build the tile image index (sth like a lookup from color to hash)
  Table: (Hash could also be replaced by Instagram object id (subject to be verfied))
  |Color|Hash|filename|
-
