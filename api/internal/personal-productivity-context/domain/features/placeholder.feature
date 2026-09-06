Feature: Domain placeholder
  The domain package starts with a placeholder entry point
  so domain behavior can grow here over time.

  Scenario: Creating a domain placeholder
    When I create a domain placeholder
    Then the placeholder should exist
