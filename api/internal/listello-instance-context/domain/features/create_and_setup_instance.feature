Feature: Create and set up a Listello instance
  As a user
  I want to create a Listello instance, choose how it is hosted, and initialize persistence
  So that I can finish setup and start using the workspace

  Scenario: Creating an instance
    When the user creates an instance
    Then a "InstanceCreated" event should have occurred
    And the instance should exist

  Scenario: Creating an instance assigns an ID prefixed with LI_
    When the user creates an instance
    Then the instance should have an ID prefixed with "LI_"
    And the instance ID after the prefix "LI_" should be a UUID

  Scenario Outline: Selecting a hosting mode
    Given an instance exists
    When the user selects hosting mode "<mode>"
    Then a "HostingModeSelected" event should have occurred
    And the instance should have hosting mode "<mode>"

    Examples:
      | mode           |
      | local          |
      | standalone-web |

  Scenario: Selecting an unknown hosting mode fails
    Given an instance exists
    When the user selects hosting mode "cloud"
    Then selecting the hosting mode should fail with error "hosting mode is not supported"

  Scenario: Selecting a persistence location for a local instance
    Given an instance exists
    And the user selects hosting mode "local"
    And the parent directory of "/var/listello" exists
    And the parent directory of "/var/listello" is writable
    And the directory "/var/listello" does not exist
    When the user selects persistence location "/var/listello"
    Then a "LocalPersistenceLocationSelected" event should have occurred
    And the instance should have persistence location "/var/listello"

  Scenario: Selecting a persistence location fails when the parent directory does not exist
    Given an instance exists
    And the user selects hosting mode "local"
    And the parent directory of "/var/listello" does not exist
    When the user selects persistence location "/var/listello"
    Then selecting the persistence location should fail with error "parent directory of persistence location does not exist"

  Scenario: Selecting a persistence location fails when the parent directory is not writable
    Given an instance exists
    And the user selects hosting mode "local"
    And the parent directory of "/var/listello" exists
    And the parent directory of "/var/listello" is not writable
    When the user selects persistence location "/var/listello"
    Then selecting the persistence location should fail with error "parent directory of persistence location is not writable"

  Scenario: Selecting a persistence location fails when the directory already exists
    Given an instance exists
    And the user selects hosting mode "local"
    And the parent directory of "/var/listello" exists
    And the parent directory of "/var/listello" is writable
    And the directory "/var/listello" exists
    When the user selects persistence location "/var/listello"
    Then selecting the persistence location should fail with error "persistence location already exists"

  Scenario: Selecting a persistence location for a standalone-web instance fails
    Given an instance exists
    And the user selects hosting mode "standalone-web"
    When the user selects persistence location "/var/listello"
    Then selecting the persistence location should fail with error "persistence location is only applicable for local"

  Scenario: Initializing persistence
    Given an instance exists
    And the user selects hosting mode "local"
    And the parent directory of "/var/listello" exists
    And the parent directory of "/var/listello" is writable
    And the directory "/var/listello" does not exist
    And the user selects persistence location "/var/listello"
    When the user initializes persistence
    Then a "PersistenceInitialized" event should have occurred with hosting mode of "local" and persistence location "/var/listello"
    And the instance should have persistence initialized
    And the instance should have the local database initialized

  Scenario: Completing setup
    Given an instance exists
    And the user selects hosting mode "local"
    And the parent directory of "/var/listello" exists
    And the parent directory of "/var/listello" is writable
    And the directory "/var/listello" does not exist
    And the user selects persistence location "/var/listello"
    And the user initializes persistence
    When the user completes setup
    Then a "SetupCompleted" event should have occurred
    And the instance should have setup completed

  Scenario: Creating a local instance through to setup completed
    When the user creates an instance
    And the user selects hosting mode "local"
    And the parent directory of "/var/listello" exists
    And the parent directory of "/var/listello" is writable
    And the directory "/var/listello" does not exist
    And the user selects persistence location "/var/listello"
    And the user initializes persistence
    And the user completes setup
    Then a "InstanceCreated" event should have occurred
    And a "HostingModeSelected" event should have occurred
    And a "LocalPersistenceLocationSelected" event should have occurred
    And a "PersistenceInitialized" event should have occurred with hosting mode of "local" and persistence location "/var/listello"
    And a "SetupCompleted" event should have occurred
    And the instance should have hosting mode "local"
    And the instance should have persistence location "/var/listello"
    And the instance should have persistence initialized
    And the instance should have the local database initialized
    And the instance should have setup completed
