# TODO
1. [ ] (HIGH) Update vulnerable packages
2. [ ] (HIGH) Add github workflows back (especially for pr tests)
3. [ ] (MED) Rename module (metafates/mangal => amonull/rengal)
4. [ ] (LOW) Update README to describe what rengal is and why it exists
5. [ ] (LOW) Investigate [plugin (golang)](https://pkg.go.dev/plugin) or [plugin (hashicorp)](https://github.com/hashicorp/go-plugin) support (windows does not seem to support this so may need to keep supporting lua scrapers as a backup)

## TODO.5 (if 5 can be done)
1. [ ] (HIGH) define plugin writing documentation steps (inputs and outputs should be clearly defined before any impl is started)
2. [ ] (MED) replacing custom lua scrapers with golang
3. [ ] (LOW) remove anilist package from repo in favour of using  

### TODO.5.2 (if 5.2 can be done)
- [ ] (HIGH) remove builtin scrapers in favour of plugin scrapers on a separate repo
- [ ] (HIGH) allow loading plugins by their module name (i.e. github.com/amonull/rengal)
- [ ] (MED) push scraper repos onto pkg.go.dev
- [ ] (LOW) have a repo for custom user scrapers for anyone to add their scraper onto