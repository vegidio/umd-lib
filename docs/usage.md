# Usage

With the library properly [installed](installation.md), you just need to initialize **UMD** by using `umd.New` passing any metadata or callback (if needed, otherwise just pass `nil`). Then, with the newly `Umd` object you can call the method `FindExtractor` passing that URL that you want to query.

**UMD** will automatically detect what site/content you're trying to fetch media information; if the URL is not supported then it will return an error.

If everything goes well and **UMD** detects the URL returning a suitable extractor, you can use the methods below:

## QueryMedia()

```go linenums="1"
u := umd.New(nil, nil)
extractor, _ := u.FindExtractor("https://www.reddit.com/user/atomicbrunette18")
resp, _ := extractor.QueryMedia(100, make([]string, 0), true)
```
