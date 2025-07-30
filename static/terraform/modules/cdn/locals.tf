locals {
  origin_id = uuid()

  # AWS Managed Caching Policy (CachingDisabled)
  # https://docs.aws.amazon.com/AmazonCloudFront/latest/DeveloperGuide/using-managed-cache-policies.html
  caching_policy_caching_disabled = "4135ea2d-6df8-44a3-9df3-4b5a84be39ad"

  # AWS Managed Origin Request Policy (AllViewer)
  # https://docs.aws.amazon.com/AmazonCloudFront/latest/DeveloperGuide/using-managed-origin-request-policies.html
  origin_request_policy_all_viewer = "216adef6-5c7f-47e4-b989-5492eafa07d3"
}
