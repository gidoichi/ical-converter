[![Publish](https://github.com/gidoichi/ical-converter/actions/workflows/publish.yml/badge.svg)](https://hub.docker.com/repository/docker/gidoichi/ical-converter/general)

# ical-converter

ical-converter convertes iCalendar (RFC5545) components to event component. It is running as a server to get iCalendar from URL and convert to VEVENT contained iCalendar to register external calendar service.

Some component specified properties are also converted to corresponding properties.

For example, given iCalendar
```
BEGIN:VCALENDAR
PRODID:-//xyz Corp//NONSGML PDA Calendar Version 1.0//EN
VERSION:2.0
BEGIN:VTODO
DTSTAMP:19960704T120000Z
UID:uid1@example.com
DTSTART:19960918T143000Z
DUE:19960920T220000Z
SUMMARY:Networld+Interop Conference
END:VTODO
END:VCALENDAR
```

is converted to
```
BEGIN:VCALENDAR
PRODID:-//xyz Corp//NONSGML PDA Calendar Version 1.0//EN
VERSION:2.0
BEGIN:VEVENT
DTSTAMP:19960704T120000Z
UID:uid1@example.com
DTSTART:19960918T143000Z
DTEND:19960920T220000Z
SUMMARY:Networld+Interop Conference
END:VEVENT
END:VCALENDAR
```

In this case, todo component is converted to event conponent. And because of due property (DUE) does not exist at vevent property, it converted to date-time end property (DTEND).

## Usage

```console
$ export PORT='8080'
$ export ICAL_CONVERTER_ICS_URL='https://example.com/remote/ical.ics'
$ go run ./main.go &
$ curl localhost:8080
BEGIN:VCALENDAR
...
END:VCALENDAR
```

## Configuration
### Environment variables
| NAME                   | DESCRIPTION                   |
|------------------------|-------------------------------|
| PORT                   | Lisning requests on this port |
| ICAL_CONVERTER_ICS_URL | Remote ical file server       |

## Detailed conversion rule
iCalendar is converted following three times.

### Non-Standard Properties
| SERVICE                        | PARSER                           |
|--------------------------------|----------------------------------|
| [2Do](https://www.2doapp.com/) | [two_do](/infrastructure/two_do) |

### Scheduled components
[converter](/usecase/converter.go)

### Filtering components
[service](/application/service.go)
