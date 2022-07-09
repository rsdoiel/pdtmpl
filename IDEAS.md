
Ideas
=====

Questions
---------

- What would be the minimum additional tools besides GNU Make and Pandoc needed to make it easy to maintain my website?
    - blogit for generating JSON presentation of blog post and rendering to the appropriate location of blog
    - pdtmpl for rendering
    - sitemapper for generating JSON sitemap files in easy directory
    - feeder for generating reallysimple feeds and RSS 2.0 feeds in each directory


Improvements
------------

- a "-walker" option that would walk the file system looking for JSON document to turn into Markdown then into HTML so that the blog is always in a data format
- use an in memory file system for temp files, this could speed up processing and also avoid the messy issue of file permissions of the content that gets buffered to disk and passed via "--template" and "--metadata-file".


