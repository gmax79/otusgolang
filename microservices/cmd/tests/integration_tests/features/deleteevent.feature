# file: features/deleteevent.feature

Feature: Delete event
    Using mycalendar api
    I can delete event by date

Scenario: Add event for delete at next scenario
    When I create event at "2020-12-31 23:00:00" with "Happy New Year"
    Then response code should be 200

Scenario: Delete event
    When I delete event "Happy New Year" at "2020-12-31 23:00:00"
    Then response code should be 200
