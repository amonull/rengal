# TODO
1. [ ] Investigate [plugin (golang)](https://pkg.go.dev/plugin) or [plugin (hashicorp)](https://github.com/hashicorp/go-plugin) support (windows does not seem to support this so may need to keep supporting lua scrapers as a backup)
2. [ ] Update vulnerable packages
3. [ ] Rename module (metafates/mangal => amonull/rengal)
4. [ ] Update README to describe what rengal is and why it exists
5. [ ] Add github workflows back (especially for pr tests)

## TODO.1 (if 1 can be done)
1. [ ] define plugin writing documentation steps (inputs and outputs should be clearly defined before any impl is started)
2. [ ] replacing custom lua scrapers with golang
3. [ ] remove anilist package from repo in favour of using  

### TODO.1.2 (if 1.2 can be done)
- [ ] remove builtin scrapers in favour of plugin scrapers on a separate repo
- [ ] have a repo for custom user scrapers for anyone to add their scraper onto
- [ ] push scraper repos onto pkg.go.dev
- [ ] allow loading plugins by their module name (i.e. github.com/amonull/rengal)