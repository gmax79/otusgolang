# file: featurer/check.feature

Feature: Add new event
    Using mycalendar api
    I can add new event

Scenario: Add event
    When I create event at "2020-10-22 18:00:00" with "Maks birthday"
    Then response code should be 200

Scenario: Get event
    When I get event at "2020-10-22 18:00:00"
    Then response should be "Maks birthday"
