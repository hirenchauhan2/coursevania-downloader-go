# Coursevania[dot]com course downloader

This repo is just to try my skills in Golang and not to spread the copyright material of course creators. I wanted to learn the rest APIs and file downloading and this popped up in my phone's notification from google about this site.

When I checked the site, I did not had a proper navigation to course and had to struggle a lot to visit the lectures. Hence, I did the backtrack of network calls and from there I decided to go for this approach to get all the files in one-shot.

## Usage of the CLI app

```
coursevaniadownloader -course="Name of the course"

# eg.
# coursevaniadownloader -course="The Complete 2020 Flutter Development Bootcamp with Dart"
```

It will create the course folder and download the course in it in the working directory only.

Please note that when you visit the site you will get the string prefixed like `"[coursevania[dot]com] Course Name"`. You don't need to pass this prefix, its handled internally.

Again, if this hurts someones feeling, then please create an issue and I will straight away delete this repo.

Thanks.
