FROM node:22-alpine
WORKDIR /app
COPY web/package.json web/package-lock.json* ./
RUN npm install
COPY web/ .
ENV NEXT_TELEMETRY_DISABLED=1
EXPOSE 3000
CMD ["npm", "run", "dev"]
