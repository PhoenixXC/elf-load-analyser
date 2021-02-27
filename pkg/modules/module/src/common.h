#ifndef XC_COMMON_H
#define XC_COMMON_H

struct _common {
    u64 ts;
};

#define TDATA(struct_name, args) \
    typedef struct struct_name { \
        u64 ts;             \
        args;                    \
    } __attribute__((packed)) struct_name;

static void* init_tdata(void *data) {
    u64 ns = bpf_ktime_get_ns();
    struct _common *tmp = (struct _common *)data;                                  
    bpf_probe_read_kernel(&tmp->ts, sizeof(u64), (void *)(&ns)); 
    return data;
}

#endif