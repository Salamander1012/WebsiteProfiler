# Website Profiler

## Requirements

This project requires Go. [See Go installation guide here](https://golang.org/doc/install)

## Build and Run

To build this project run 
```
make build
```

To profile a website 
```
./sitestat --url <url_string> --profile <profile_size>
```

You can also see usage instructions
```
./sitestat --help
```

## Results

Cloudflare is fast!

### Google

![google](https://github.com/Salamander1012/WebsiteProfiler/blob/main/results/google.png)

### Youtube

![youtube](https://github.com/Salamander1012/WebsiteProfiler/blob/main/results/youtube.png)

### Amazon

![amazon](https://github.com/Salamander1012/WebsiteProfiler/blob/main/results/amazon.png)

### My Cloudflare Worker

Profile:

![my_worker_profile](https://github.com/Salamander1012/WebsiteProfiler/blob/main/results/my_worker_profile.png)

Response:

![my_worker_response](https://github.com/Salamander1012/WebsiteProfiler/blob/main/results/my_worker_response.png)