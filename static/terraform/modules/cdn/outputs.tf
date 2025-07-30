output "cloudfront_url" {
  value       = "https://${aws_cloudfront_distribution.this.domain_name}"
  description = "CloudFront distribution URL"
}
