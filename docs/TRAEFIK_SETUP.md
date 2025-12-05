# Traefik Configuration for Dokploy

This project is configured to work with Traefik reverse proxy in Dokploy.

## Required Environment Variable

Add this to your `.env` file or Dokploy environment variables:

```bash
# Traefik External Domain
BACKEND_EXTERNAL_DOMAIN=api.yourdomain.com
```

## How It Works

The `docker-compose.prod.yml` includes Traefik labels that:

1. **Connect to Dokploy's Traefik network**
   - Joins `dokploy-network` (external network managed by Dokploy)
   - Maintains internal network for postgres/grafana communication

2. **Configure HTTPS with Let's Encrypt**
   - Automatic SSL certificate from Let's Encrypt
   - Routes traffic based on hostname

3. **Load Balancer Configuration**
   - Routes external traffic to your app's internal port (defined by `APP_PORT`)

## Traefik Labels Explained

```yaml
labels:
  # Enable Traefik for this service
  - "traefik.enable=true"
  
  # Which network Traefik should use to reach this container
  - "traefik.docker.network=dokploy-network"
  
  # Routing rule: match requests to your domain
  - "traefik.http.routers.mynute-backend-prod-http.rule=Host(`${BACKEND_EXTERNAL_DOMAIN}`)"
  
  # Use HTTPS (websecure entrypoint)
  - "traefik.http.routers.mynute-backend-prod-http.entrypoints=websecure"
  
  # Get SSL certificate from Let's Encrypt
  - "traefik.http.routers.mynute-backend-prod-http.tls.certresolver=letsencrypt"
  
  # Tell Traefik which port your app listens on
  - "traefik.http.services.mynute-backend-prod-http.loadbalancer.server.port=${APP_PORT}"
```

## Setup in Dokploy

1. **Set the domain in environment variables:**
   ```bash
   BACKEND_EXTERNAL_DOMAIN=api.yourdomain.com
   APP_PORT=4000
   ```

2. **Configure DNS:**
   - Point `api.yourdomain.com` to your Dokploy server IP
   - Wait for DNS propagation

3. **Deploy:**
   - Dokploy will automatically configure Traefik
   - Let's Encrypt will issue SSL certificate
   - Your API will be accessible at `https://api.yourdomain.com`

## Benefits

✅ **Automatic HTTPS** - Let's Encrypt certificates  
✅ **No port exposure** - No need for `ports:` mapping  
✅ **Multiple domains** - Can run multiple apps on same server  
✅ **Load balancing** - Built-in load balancing if scaling  
✅ **Automatic renewal** - SSL certificates auto-renew  

## Troubleshooting

### Service not accessible via domain

**Check Traefik is running:**
```bash
docker ps | grep traefik
```

**Check container is on dokploy-network:**
```bash
docker inspect <container-name> | grep dokploy-network
```

**View Traefik logs:**
```bash
docker logs dokploy-traefik
```

**Verify DNS is pointing to server:**
```bash
dig api.yourdomain.com
```

### SSL certificate not issued

- Wait a few minutes for Let's Encrypt validation
- Ensure port 80 and 443 are open on your server
- Check domain DNS is correct
- View Traefik logs for certificate errors

## Related Documentation

- [Dokploy Deployment Guide](./DOKPLOY_DEPLOYMENT.md)
- [Traefik Official Docs](https://doc.traefik.io/traefik/)
- [Let's Encrypt](https://letsencrypt.org/)
