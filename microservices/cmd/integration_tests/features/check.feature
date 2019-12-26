# file: featurer/check.feature

Feature: new event
    In order to use mycalendar api
    As an API calendar
    I need add new event

Scenario: Add event
    When I send "GET" requst to "/users"
    Then response code should be 200
