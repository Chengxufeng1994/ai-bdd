Feature: Service version

  So that a deployed instance can be identified after the fact, a running
  service reports the version it was built from.

  This feature carries no domain meaning. It exists to keep one thin path
  through the whole chain alive — OpenAPI contract, generated server, handler,
  scenario, godog — so that a break anywhere in it is caught by a test rather
  than discovered while writing the first real feature.

  Rule: A running service reports the version it was built from

    @smoke
    Scenario: Asking a running service for its version
      Given the service is running at version "1.2.3"
      When a client asks for the service version
      Then the service reports version "1.2.3"
