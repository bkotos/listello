Feature: Create and set up a Listello instance
  As a user
  I want to create a Listello instance, choose how it is hosted, and initialize persistence
  So that I can finish setup and start using the workspace

  Scenario: Creating an instance
    When the user creates an instance
    Then a "InstanceCreated" event should have occurred
    And the instance should exist

  Scenario Outline: Selecting a hosting mode
    Given an instance exists
    When the user selects hosting mode "<mode>"
    Then a "HostingModeSelected" event should have occurred
    And the instance should have hosting mode "<mode>"

    Examples:
      | mode           |
      | local          |
      | standalone-web |

  @wip
  Scenario: Selecting a persistence location for a local instance
    Given an instance exists
    And the user selects hosting mode "local"
    When the user selects persistence location "/var/listello"
    Then a "LocalPersistenceLocationSelected" event should have occurred
    And the instance should have persistence location "/var/listello"

  @wip
  Scenario: Selecting a persistence location for a standalone-web instance fails
    Given an instance exists
    And the user selects hosting mode "standalone-web"
    When the user selects persistence location "/var/listello"
    Then selecting the persistence location should fail with error "persistence location is only applicable for local"

  @wip
  Scenario: Initializing persistence
    Given an instance exists
    And the user selects hosting mode "local"
    And the user selects persistence location "/var/listello"
    When the user initializes persistence
    Then a "PersistenceInitialized" event should have occurred
    And the instance should have persistence initialized

  @wip
  Scenario: Completing setup
    Given an instance exists
    And the user selects hosting mode "local"
    And the user selects persistence location "/var/listello"
    And the user initializes persistence
    When the user completes setup
    Then a "SetupCompleted" event should have occurred
    And the instance should have setup completed

  @wip
  Scenario: Creating a local instance through to setup completed
    When the user creates an instance
    And the user selects hosting mode "local"
    And the user selects persistence location "/var/listello"
    And the user initializes persistence
    And the user completes setup
    Then a "InstanceCreated" event should have occurred
    And a "HostingModeSelected" event should have occurred
    And a "LocalPersistenceLocationSelected" event should have occurred
    And a "PersistenceInitialized" event should have occurred
    And a "SetupCompleted" event should have occurred
    And the instance should have hosting mode "local"
    And the instance should have persistence location "/var/listello"
    And the instance should have persistence initialized
    And the instance should have setup completed
