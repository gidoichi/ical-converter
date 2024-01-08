[![Publish](https://github.com/gidoichi/ical-converter/actions/workflows/publish.yml/badge.svg)](https://hub.docker.com/repository/docker/gidoichi/ical-converter/general)

# ical-converter

ical-converter convertes iCalendar (RFC5545) components to event component.

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
