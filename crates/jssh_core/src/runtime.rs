use crate::bundle::{NativeModule, SCRIPT_MODULE};
use anyhow::{anyhow, Result};
use rquickjs::{
    loader::{
        BuiltinLoader, BuiltinResolver, FileResolver, ModuleLoader, NativeLoader, ScriptLoader,
    },
    CatchResultExt, Context, Ctx, Error, Function, Module, Promise, Runtime, Value,
    qjs
};

struct RuntimeCore {
    rt: Runtime,
    ctx: Context,
}

impl RuntimeCore {
    pub fn new() -> Result<Self> {
        let rt = Runtime::new()?;
        let ctx = Context::full(&rt)?;

        let resolver = (
            BuiltinResolver::default()
                .with_module("bundle/script_module")
                .with_module("bundle/native_module"),
            FileResolver::default()
                .with_path("./")
                .with_path("../../target/debug")
                .with_native(),
        );
        let loader = (
            BuiltinLoader::default().with_module("bundle/script_module", SCRIPT_MODULE),
            ModuleLoader::default().with_module("bundle/native_module", NativeModule),
            ScriptLoader::default(),
            NativeLoader::default(),
        );
        rt.set_loader(resolver, loader);

        Ok(RuntimeCore { rt, ctx })
    }

    pub fn eval<F>(&self, name: &str, source: &str, f: F) -> Result<()>
    where
        F: FnOnce(Ctx, Promise) -> Result<()> + Send,
    {
        self.ctx.with(|ctx| {
            let ret = Module::evaluate(ctx.clone(), name, source)?;
            f(ctx, ret)
        })
    }

    pub fn eval_anonymous<F>(&self, source: &str, f: F) -> Result<()>
    where
        F: FnOnce(Ctx, Promise) -> Result<()> + Send,
    {
        self.eval("anonymous", source, f)
    }
}

#[cfg(test)]
mod tests {
    use crate::runtime::*;
    use anyhow::Result;

    #[test]
    fn it_works() -> Result<()> {
        let rt = RuntimeCore::new()?;
        match rt.eval_anonymous("throw new Error(123)", |ctx, p| {
            match p.finish::<Value>() {
                Ok(result) => {
                    println!("result: {:?}", result);
                }
                Err(err) => {
                    println!("error: {}", err);
                }
            }
            Ok(())
        }) {
            Ok(_) => {
                println!("done")
            }
            Err(e) => {
                eprintln!("eval error: {}", e)
            }
        };
        Ok(())
    }
}
