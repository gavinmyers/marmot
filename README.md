marmot ~ Go, Redis & GitHub powered blog.
======

Marmot is bare-bones blogging application that listens for GitHub commits and caches them locally in a Redis database.

That's it!

For example, here's my blog:
https://github.com/gavinmyers/blog

If you run Marmot and configure it to point to that blog you can access the pages like this:
http://localhost/theme/index.html

# Installation (ideal)
- Install the required software (go, redis)
- Clone the marmot repository (the repo you are reading right now)
- Create a new repository for your blog (like this https://github.com/gavinmyers/blog)
- Add a webhook url pointing to your marmot instance
- Access your blog and let everything configure itself (this doesn't happen)

Right now it isn't that simple, it easily can be, but you can see in the source code much of what would be automated is hardcoded to my own blog. 
