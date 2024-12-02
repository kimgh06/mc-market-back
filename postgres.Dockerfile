FROM postgres:16

# Install dependencies
RUN apt update
RUN apt install -y wget g++ make git postgresql-server-dev-16

# Download and install mecab_ko
RUN wget -O mecab_ko.tar.gz https://bitbucket.org/eunjeon/mecab-ko/downloads/mecab-0.996-ko-0.9.2.tar.gz
RUN mkdir mecab_ko
RUN tar zxfv mecab_ko.tar.gz -C mecab_ko --strip-components=1
WORKDIR mecab_ko
RUN ./configure
RUN make
RUN make install
WORKDIR /

RUN ldconfig

# Download and install mecab_ko_dic
RUN wget -O mecab_ko_dic.tar.gz https://bitbucket.org/eunjeon/mecab-ko-dic/downloads/mecab-ko-dic-2.1.1-20180720.tar.gz
RUN mkdir mecab_ko_dic
RUN tar zxfv mecab_ko_dic.tar.gz -C mecab_ko_dic --strip-components=1
WORKDIR mecab_ko_dic
RUN ./configure
RUN make
RUN make install
WORKDIR /

# Download and install textsearch_ko PostgreSQL extension
RUN git clone https://github.com/i0seph/textsearch_ko
WORKDIR textsearch_ko
RUN make
RUN make install
#RUN psql -U postgres -f ts_mecab_ko.sql
#RUN echo 'psql -U postgres -d $POSTGRES_DB -f /textsearch_ko/ts_mecab_ko.sql' > ./init_extensions.sh
#COPY ./postgresql_initializer.sh /docker-entrypoint-initdb.d/
#RUN chmod 755 /docker-entrypoint-initdb.d/postgresql_initializer.sh
RUN cp ./ts_mecab_ko.sql /docker-entrypoint-initdb.d
WORKDIR /

