# AggreGator

This application can be used to scrape RSS-feeds and store the fetched posts into a postgres database.

To use this application, postgres and go are required. Given you have them, you can install the binary using `go install`:

    go install https://github.com/KrysPow/Blog_Aggregator

In order for the application to work, you need a configuration file in your home directory `~/.gatorconfig.json` which holds the url to your database. Mine looks like this:

    {
     "db_url": "postgres://postgres:postgres@localhost:5432/gator?sslmode=disable"
    }

Now you can use the applications. Here are a few commands:

- `register` user_name
- `login` user_name
- `addfeed` feed_name feed_url
- `agg` time_interval
- `browse` number_of_shown_posts