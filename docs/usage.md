# Usage

With the library properly [installed](installation.md), you just need to initialize create an `UMD` type with `umd.New` passing the URL of the website that you want to fetch media information. **UMD-lib** will automatically detect what site/content you're trying to fetch media information; if the URL is not supported then it will return an error.

If everything goes well and **UMD-lib** detects the URL, you can use the methods below:

## QueryMedia()

```go linenums="1"
u, err := umd.New("https://www.reddit.com/user/atomicbrunette18", nil, nil)
resp, err := u.QueryMedia(100, make([]string, 0), true)
```
