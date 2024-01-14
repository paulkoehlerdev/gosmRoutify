# gosmRoutify

> [!CAUTION]
> This project is a work in progress.

> [!IMPORTANT]
> This project does not try to be a full-fledged routing engine.
> It's a simple development project to learn about routing algorithms and data structures. \
> It's not meant to be used in production. Please don't try to use it in production.
> 
> Even thought this project is related to a course at the University of Applied Sciences in Munich,
> it is not an official project of the university. It does not represent the opinion of the university
> or any of its employees.

## Demo

There is a public demo available at [https://demo.gosmroutify.xyz](https://demo.gosmroutify.xyz).
It's a simple web interface to test the routing engine. If there are any bugs or problems, please open an issue.
This demo may be offline or unusable from time to time, because I'm using it for testing and development.

## Documentation

The full documentation is only available in German. \
You can find it [here (https://docs.gosmroutify.xyz)](https://docs.gosmroutify.xyz).

## Special Thanks

Thanks to the awesome OpenSource community around OSM for their great work on the OSM dataset.
Special thanks to the owners and maintainers of following libraries, which I either used or was strongly inspired by:

- [`graphhopper` https://github.com/graphhopper/graphhopper](https://github.com/graphhopper/graphhopper)
- [`valhalla` https://github.com/valhalla/valhalla](https://github.com/valhalla/valhalla)

I also want to thank [Geofabrik GmbH](https://www.geofabrik.de/) for providing the OSM data extracts for testing and experiencing this project.

## Used Software and Libraries

### Backend

- [`Go` https://go.dev](https://go.dev/)
- [`go-sqlite3` https://github.com/mattn/go-sqlite3](https://github.com/mattn/go-sqlite3)
- [`SQLite` https://www.sqlite.org](https://www.sqlite.org/index.html)

### Frontend

- [`Vue.js` https://vuejs.org](https://vuejs.org/)
- [`Leaflet` https://leafletjs.com](https://leafletjs.com/)
- [`Bootstrap` https://getbootstrap.com](https://getbootstrap.com/)
- [`Bootstrap Icons` https://icons.getbootstrap.com](https://icons.getbootstrap.com/)
- [`lodash` https://lodash.com](https://lodash.com/)
- [`mitt` https://github.com/developit/mitt](https://github.com/developit/mitt)

### Documentation

- [`Docusaurus` https://docusaurus.io](https://docusaurus.io/)

## Hosting

The demo backend is currently hosted on a VPS and protected by [Cloudflare](https://www.cloudflare.com/).
The demo frontend and documentation is hosted on [Vercel](https://vercel.com).