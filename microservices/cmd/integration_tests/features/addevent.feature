# file: features/addevent.feature

Feature: Add new event
    Using mycalendar api
    I can add new event
    And when get event back

Scenario: Add event
    When I create event at "2020-10-22 18:00:00" with "Maks birthday"
    Then response code should be 200

Scenario: Calculate events at day
    When I get events at "2020-10-22"
    Then count should be 1
