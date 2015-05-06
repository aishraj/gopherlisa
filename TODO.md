## Things to do

- Figure out the new workflow
- Add CSS
- (Resize) and store tile images. (Each file would be identified by its hash/id ) and also the target image.
- Add a DB connection to mysql (needed for storing index) or to redis (depends)
- Build the tile image index (sth like a lookup from color to hash)
  Table: (Hash could also be replaced by Instagram object id (subject to be verfied))
  |Color|Hash|filename|
-
Wokflow:
-> Landing page with Login button (if not logged in)
-> Upload Page
-> Search Page (like google with the name Hi Blah!)
-> (Progress Bar)
-> Catalogue (compute and index images)
-> Use the algo
-> Show it to the user.
