# file: features/calcevent.feature

Feature: Add events and then calculate them
    Using mycalendar api
    I can add events
    And when calculate them

Scenario: Add event 1
    When I create event at "2021-01-04 18:00:00" with "New Year Party"
    Then response code should be 200

Scenario: Add event 2
    When I create event at "2021-01-04 12:00:00" with "Skiing"
    Then response code should be 200

Scenario: Add event 3
    When I create event at "2021-01-07 16:00:00" with "Cristmas"
    Then response code should be 200

Scenario: Add event 4
    When I create event at "2021-01-11 12:00:00" with "Maks birthday"
    Then response code should be 200

Scenario: Calculate events at day
    When I get events at "2021-01-04"
    Then count should be 2

Scenario: Calculate events at week 1
    When I get week events at "2021-01"
    Then count should be 3

Scenario: Calculate events at week 2
    When I get week events at "2021-02"
    Then count should be 1

Scenario: Calculate events at week 3
    When I get week events at "2021-03"
    Then count should be 0

Scenario: Calculate events at month
    When I get month events at "2021-01"
    Then count should be 4
