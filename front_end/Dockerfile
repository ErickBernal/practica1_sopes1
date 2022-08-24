FROM node:alpine3.16
WORKDIR /app
COPY . .
RUN npm install
RUN npm install axios
RUN npm install @fortawesome/react-fontawesome
RUN npm install @fortawesome/free-solid-svg-icons
RUN npm install reactstrap  

EXPOSE 3000
CMD ["npm", "start"]