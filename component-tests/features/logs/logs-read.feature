Feature: Log Read Tests

  Scenario: Read Logs
    Given I reset DB and mocks
    When I read logs for resource type "account" and store "00000000-0000-0000-0000-000000000000"
    Then I see 200 status code in response

