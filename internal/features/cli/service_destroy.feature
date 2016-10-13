@service @service_destroy
Feature: Service destroy

  Scenario: Non logged service destroy
    Given I setup ernest with target "https://ernest.local"
    And I logout
    When I run ernest with "service destroy"
    Then The output should contain "You're not allowed to perform this action, please log in"
    When I run ernest with "service destroy destroyable"
    Then The output should contain "You're not allowed to perform this action, please log in"

  Scenario: Logged service destroy unexisting
    Given I setup ernest with target "https://ernest.local"
    And I'm logged in as "usr" / "pwd"
    When I run ernest with "service destroy"
    Then The output should contain "You should specify an existing service name"
    When I run ernest with "service destroy unexisting --force"
    Then The output should contain "Specified service name does not exist"

  Scenario: Logged service destroy
    Given I setup ernest with target "https://ernest.local"
    And I'm logged in as "usr" / "pwd"
    And the datacenter "test_dc" does not exist
    And I run ernest with "datacenter create aws --token tmp_token --secret tmp_secret --region tmp_region --fake test_dc"
    And The service "destroyable" does not exist
    And I run ernest with "service apply internal/definitions/destroyable.yml"
    When I run ernest with "service destroy destroyable --force"
    Then The output line number "0" should contain "Starting environment deletion"

