# Stargazer

Stargazer is a collection of controllers for Nimbess. Each controller monitors
Kubernetes resources using the Kubernetes API and performs different tasks in
response to received events from the API. One of these watched resources includes
the Unified Network Policy (UNP) Custom Resource Definition (CRD) which is defined
in this project. The UNP is a core component of Nimbess, which uses policy to
manage container network configuration and intent.

For more information on how to deploy or use Stargazer within Nimbess please
refer to the docs repo here:

[www.github.com/nimbess/nimbess-specs]
