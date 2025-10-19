resource "aws_s3_bucket_public_access_block" "good_example" {
  bucket              = aws_s3_bucket.example.id
  block_public_policy = true
}