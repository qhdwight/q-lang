#pragma once

#include <map>

#include <parser/node/implementation/integer_node.hpp>
#include <parser/node/structure/parse_with_descriptor_node.hpp>

namespace ql::parser {
    class ImplementFunctionNode : public ParseWithDescriptorNode {
    protected:
        using VariableNodeFactory = std::function<std::shared_ptr<ParseNode>(std::string&&, Tokens&&, AbstractNode::ParentRef)>;

        std::map<std::string, VariableNodeFactory> m_NameToVariableNode;

        template<typename TNode>
        void registerVariableNode(std::string_view nodeName) {
            // TODO use forwarding?
            m_NameToVariableNode.emplace(nodeName, [](auto&& block, auto&& tokens, auto parent) {
                auto node = std::make_shared<TNode>(std::forward<decltype(block)>(block), std::forward<decltype(tokens)>(tokens), parent);
                node->parse();
                return node;
            });
        }
    public:
        ImplementFunctionNode(std::string&& body, std::string_view const& innerBody, Tokens&& tokens, ParentRef const& parent);

        void parse() override;
    };
}
