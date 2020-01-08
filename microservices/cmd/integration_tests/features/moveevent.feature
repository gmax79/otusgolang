# file: features/moveevent.feature

Feature: Move event
    Using mycalendar api
    I can move event from one date to another

Scenario: Add event for move at next scenario
    When I create event at "2020-09-01 10:00:00" with "Welcome to school"
    Then response code should be 200

Scenario: Delete event if it exists
    When I delete event "Welcome to school" at "2020-09-01 11:30:00"
    Then response code should be 200

Scenario: Move event
    When I move event "Welcome to school" at "2020-09-01 10:00:00" to "2020-09-01 11:30:00"
    Then response code should be 200

Scenario: Calculate events at day
    When I get events at "2020-09-01"
    Then count should be 1
