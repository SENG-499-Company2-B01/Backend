FROM python:3.10.6-slim-buster 

RUN pip3 install pymongo pandas python-dotenv bcrypt

COPY . . 

CMD ["python3","db_seed.py"]