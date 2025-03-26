# Security Policy

## Reporting a Vulnerability

We take security vulnerabilities seriously. We appreciate your efforts to responsibly disclose your findings.

To report a security issue, please submit a security issue to this repository with a description of the issue, the steps you took to create it, affected versions, and if known, mitigations. Please do not open GitHub issues for security vulnerabilities.

We'll acknowledge your submission within 48 hours and provide a detailed response within 72 hours indicating the next steps in handling your report.

## Supported Versions

| Version | Supported          |
| ------- | ------------------ |
| 1.0.x   | :white_check_mark: |
| < 1.0.0 | :x:                |

## Security Considerations

### AWS IAM Permissions

Beacon requires the `events:PutEvents` IAM permission to function. When configuring this permission:

- Follow the principle of least privilege: only grant the necessary permissions to the EC2 instance role
- Restrict the permission to specific AWS regions if not used globally
- Consider using condition keys to further restrict access, such as source IP or AWS tags

### Authentication Considerations

- Beacon uses the AWS SDK which follows AWS credential resolution chain
- Do not hardcode AWS credentials in user data scripts or environment variables
- Use IAM roles for EC2 instances when possible

### Data Considerations

- Beacon sends event data to CloudWatch Events with the following fields:
  - Instance ID
  - Project name
  - Status message
  - Custom message
- Avoid including sensitive data in event messages (passwords, tokens, etc.)
- Be aware that CloudWatch Events are accessible to anyone with appropriate AWS IAM permissions

## Dependency Management

### Third-Party Dependencies

Beacon relies on the following main dependencies:

- AWS SDK for Go v2 - For AWS API access
- Cobra - For command-line interface

All dependencies are tracked in go.mod and go.sum files. The project uses dependabot (configured in `.github/dependabot.yml`) to automatically update dependencies when security patches are available.

### Dependency Auditing

The project undergoes the following security checks:

- Weekly dependency updates via Dependabot

## Security Best Practices for Users

1. **Keep Beacon Updated**: Always use the latest version to benefit from security patches.

2. **IAM Permissions**: Configure the minimal IAM permissions required:
   ```json
   {
     "Version": "2012-10-17",
     "Statement": [
       {
         "Effect": "Allow",
         "Action": "events:PutEvents",
         "Resource": "*"
       }
     ]
   }
   ```

3. **Secure User Data**: If embedding Beacon commands in user data scripts, ensure those scripts are transmitted securely (e.g., via HTTPS) and do not contain sensitive information.

4. **Message Content**: Be cautious about the content of messages sent through Beacon. Avoid including sensitive data.

5. **CloudWatch Events Access**: Restrict access to CloudWatch Events (aka EventBridge) by implementing appropriate IAM policies.

## Release Process

Security fixes are released as soon as possible, following this process:

1. Security fixes are prioritized over feature development
2. For critical vulnerabilities, patch releases are issued within 7 days
3. Security-related releases are clearly marked in release notes

## Security Hardening

For production environments, consider these additional hardening measures:

1. Run Beacon in a container with minimal permissions
2. Consider network restrictions, allowing only necessary outbound connections
3. Implement additional logging and monitoring for Beacon operations
4. Use AWS VPC endpoints for CloudWatch to avoid data traversing the public internet

## License

This security policy and the Beacon project are under the MIT License. See the LICENSE file for more details.
