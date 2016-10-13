@datacenter @datacenter_create @datacenter_create_aws
Feature: Ernest datacenter create

  Scenario: Non logged aws datacenter creation
    Given I setup ernest with target "https://ernest.local"
    And I logout
    When I run ernest with "datacenter create aws"
    Then The output should contain "You should specify the datacenter name"
    When I run ernest with "datacenter create aws tmp_datacenter"
    Then The output should contain "You're not allowed to perform this action, please log in"

  Scenario: Logged user aws datacenter creation
    Given I setup ernest with target "https://ernest.local"
    And the datacenter "tmp_datacenter" does not exist
    And I'm logged in as "usr" / "pwd"
    When I run ernest with "datacenter create aws"
    Then The output should contain "You should specify the datacenter name"
    When I run ernest with "datacenter create aws tmp_datacenter"
    Then The output should contain "Please, fix the error shown below to continue"
    And The output should contain "- Specify a valid token with --token flag"
    And The output should contain "- Specify a valid secret with --secret flag"
    And The output should contain "- Specify a valid region with --region flag"
    When I run ernest with "datacenter create aws --token tmp_token --secret tmp_secret --region tmp_region tmp_datacenter"
    Then The output should contain "Datacenter 'tmp_datacenter' successfully created"
    When I run ernest with "datacenter list"
    Then The output should contain "tmp_datacenter"
    Then The output should contain "tmp_region"
    Then The output should contain "aws"

  Scenario: Adding an already existing aws datacenter
    Given I setup ernest with target "https://ernest.local"
    And the datacenter "tmp_datacenter" does not exist
    And I'm logged in as "usr" / "pwd"
    When I run ernest with "datacenter create aws --token tmp_token --secret tmp_secret --region tmp_region tmp_datacenter"
    And I run ernest with "datacenter create aws --token tmp_token --secret tmp_secret --region tmp_region tmp_datacenter"
    Then The output should contain "Datacenter 'tmp_datacenter' already exists, please specify a different name"


