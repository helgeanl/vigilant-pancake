
#include "driver/io.h"
#include "erl_nif.h"

static ERL_NIF_TERM io_init_nif(ErlNifEnv* env, int argc, const ERL_NIF_TERM argv[]){
    int elevator_type;
    if(!enif_get_int(env, argv[0], &elevator_type)){
        return enif_make_badarg(env);
    }
    int ret = io_init(elevator_type);
    return enif_make_int(env, ret);
}

static ERL_NIF_TERM io_set_bit_nif(ErlNifEnv* env, int argc, const ERL_NIF_TERM argv[]){
    int channel;
    if(!enif_get_int(env, argv[0], &channel)){
        return enif_make_badarg(env);
    }
    io_set_bit(channel);
    return enif_make_atom(env, "ok");
}

static ErlNifFunc nif_funcs[] = {
    {"io_init", 1, io_init_nif},
    {"io_set_bit", 1, io_set_bit_nif},
};

ERL_NIF_INIT(main, nif_funcs, NULL, NULL, NULL, NULL)

